package outputs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"bytes"

	"sync"

	"github.com/r3boot/logshipper/lib/config"
)

const (
	MAX_QUEUE_ITEMS = 32768
)

type ESShipper struct {
	client     *http.Client
	esUri      string
	todayIndex string
	batchMutex sync.RWMutex
	batchQueue [][]byte
	batchChan  chan []byte
	Control    chan int
	Done       chan bool
}

type esGetIndexData struct {
	Health     string `json:"health"`
	Status     string `json:"status"`
	Index      string `json:"index"`
	Uuid       string `json:"uuid"`
	Pri        string `json:"pri"`
	Rep        string `json:"rep"`
	NumDocs    string `json:"docs.count"`
	NumDeleted string `json:"docs.deleted"`
	Size       string `json:"store.size"`
	PriSize    string `json:"pri.store.size"`
}

type esGetIndicesData []esGetIndexData

type esCreateIndexResponseData struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index"`
}

type esCreateIndexSettingsData struct {
	NumberOfShards   int `json:"number_of_shards"`
	NumberOfReplicas int `json:"number_of_replicas"`
}

type esCreateIndexMappingsData struct {
	Properties map[string]map[string]string `json:"properties"`
}

type esCreateIndexData struct {
	Settings esCreateIndexSettingsData `json:"settings"`
	// Mappings esCreateIndexMappingsData `json:"mappings"`
}

type esAliasActionData struct {
	Index string `json:"index"`
	Alias string `json:"alias"`
}

type esAliasActionsData map[string]*esAliasActionData

type esAliasData struct {
	Actions []esAliasActionsData `json:"actions"`
}

type esBulkResponseItemShardsData struct {
	Total     int `json:"total"`
	Succesful int `json:"succesful"`
	Failed    int `json:"failed"`
}

type esBulkResponseItemData struct {
	Index       string `json:"_index"`
	Type        string `json:"_type"`
	Id          string `json:"_id"`
	Version     int    `json:"_version"`
	Result      string `json:"result"`
	Shards      esBulkResponseItemShardsData
	SeqNo       int `json:"_seq_no"`
	PrimaryTerm int `json:"_primary_term"`
	Status      int `json:"status"`
}

type esBulkResponseData struct {
	Took   int                                 `json:"took"`
	Errors bool                                `json:"errors"`
	Items  []map[string]esBulkResponseItemData `json:"items"`
}

func NewESShipper() (*ESShipper, error) {
	es := &ESShipper{
		client:    &http.Client{Transport: &http.Transport{}},
		esUri:     fmt.Sprintf("%s:%d", cfg.ES.Host, cfg.ES.Port),
		batchChan: make(chan []byte, MAX_QUEUE_ITEMS),
		Control:   make(chan int, 1),
		Done:      make(chan bool, 1),
	}

	err := es.checkIndex()
	if err != nil {
		return nil, fmt.Errorf("NewESShipper: %v", err)
	}

	go es.BatchSubmit()

	return es, nil
}

func (es *ESShipper) checkIndex() error {
	tNow := time.Now()
	es.todayIndex = fmt.Sprintf("%s-%d.%d.%d", cfg.ES.Index, tNow.Year(), tNow.Month(), tNow.Day())

	indices, err := es.getIndices()
	if err != nil {
		return fmt.Errorf("ESShipper.checkIndex: %v", err)
	}

	haveIndex := false
	for _, index := range indices {
		if index == es.todayIndex {
			haveIndex = true
			break
		}
	}

	if !haveIndex {
		_, err = es.addIndex(es.todayIndex)
		if err != nil {
			return fmt.Errorf("ESShipper.checkIndex: %v", err)
		}
		log.Debugf("ESShipper.checkIndex: created index %s", es.todayIndex)

		err = es.addAlias(es.todayIndex)
		if err != nil {
			return fmt.Errorf("ESShipper.checkIndex: %v", err)
		}
		log.Debugf("ESShipper.checkIndex: added %s to alias for %s", es.todayIndex, cfg.ES.Index)

	} else {
		log.Debugf("ESShipper.checkIndex: index %s already exists", es.todayIndex)
	}

	return nil
}

func (es *ESShipper) getIndices() ([]string, error) {
	uri := fmt.Sprintf("http://%s/_cat/indices?format=json", es.esUri)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("ESShipper.getIndices http.NewRequest: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ESShipper.getIndices client.Do: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	indices := esGetIndicesData{}
	err = decoder.Decode(&indices)
	if err != nil {
		return nil, fmt.Errorf("ESShipper.getIndices decoder.Decode: %v", err)
	}

	allIndices := []string{}
	for _, index := range indices {
		allIndices = append(allIndices, index.Index)
	}

	return allIndices, nil
}

func (es *ESShipper) addIndex(name string) (string, error) {
	idxConfig := esCreateIndexData{
		Settings: esCreateIndexSettingsData{
			NumberOfShards:   5,
			NumberOfReplicas: 1,
		},
	}

	data, err := json.Marshal(idxConfig)
	if err != nil {
		return "", fmt.Errorf("ESShipper.addIndex json.Marshal: %v", err)
	}

	uri := fmt.Sprintf("http://%s/%s", es.esUri, name)

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("ESShipper.addIndex http.NewRequest: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ESShipper.addIndex client.Do: %v", err)
	}
	defer resp.Body.Close()

	response := esCreateIndexResponseData{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	if !response.Acknowledged {
		return "", fmt.Errorf("ESShipper.addIndex: index not acknowledged")
	}

	if !response.ShardsAcknowledged {
		return "", fmt.Errorf("ESShipper.addIndex: index shards not acknowledged")
	}

	return name, nil
}

func (es *ESShipper) addAlias(name string) error {
	aliasConfig := esAliasData{
		Actions: []esAliasActionsData{},
	}

	action := esAliasActionsData{}
	action["add"] = &esAliasActionData{
		Index: name,
		Alias: cfg.ES.Index,
	}
	aliasConfig.Actions = append(aliasConfig.Actions, action)

	return nil
}

func (es *ESShipper) runBatch(batch [][]byte) error {
	esBatchRequest := ""
	for _, entry := range batch {
		esBatchRequest += fmt.Sprintf("{\"index\": {\"_index\":\"%s\",\"_type\":\"doc\"}}\n", es.todayIndex)
		esBatchRequest += fmt.Sprintf("%s\n", entry)
	}

	uri := fmt.Sprintf("http://%s/_bulk", es.esUri)

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer([]byte(esBatchRequest)))
	if err != nil {
		return fmt.Errorf("ESShipper.runBatch http.NewRequest: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return fmt.Errorf("ESShipper.runBatch client.Do: %v", err)
	}
	defer resp.Body.Close()

	response := esBulkResponseData{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return fmt.Errorf("ESShipper.runBatch decoder.Decode: %v", err)
	}

	if response.Errors {
		log.Warningf("ESShipper.runBatch: Found errors in documents")
		for _, entry := range response.Items {
			if entry["index"].Status != http.StatusCreated {
				log.Warningf("failed: %v", entry)
			}
		}
	} else {
		log.Debugf("ESShipper.runBatch: indexed %d documents in %dms", len(batch), response.Took)
	}

	return nil
}

func (es *ESShipper) BatchSubmit() {
	log.Debugf("ESShipper.BatchSubmit: Starting submit routine")
	tDuration, _ := time.ParseDuration("0s")
	for {
		time.Sleep((1 * time.Second) - tDuration)

		tStart := time.Now()

		es.batchMutex.Lock()
		close(es.batchChan)
		batch := [][]byte{}
		for entry := range es.batchChan {
			batch = append(batch, entry)
		}
		es.batchChan = make(chan []byte, MAX_QUEUE_ITEMS)
		es.batchMutex.Unlock()

		if len(batch) == 0 {
			log.Debugf("ESShipper.BatchSubmit: Queue is empty")
			continue
		}

		es.runBatch(batch)

		tDuration = time.Since(tStart)
	}
}

func (es *ESShipper) Ship(logdata chan []byte) error {

	stop_loop := false
	for {
		if stop_loop {
			break
		}

		select {
		case event := <-logdata:
			{
				es.batchMutex.Lock()
				es.batchChan <- event
				es.batchMutex.Unlock()
			}
		case cmd := <-es.Control:
			{
				switch cmd {
				case config.CMD_CLEANUP:
					{
						log.Debugf("AmqpShipper.Ship: Shutting down")
						stop_loop = true
						continue
					}
				default:
					{
						log.Warningf("AmqpShipper.Ship: Invalid command received")
					}
				}
			}
		}
	}

	es.Done <- true

	return nil
}
