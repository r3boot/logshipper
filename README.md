This is a Go-Lang based daemon, meant as a lightweight alternative to running
a JVM on all systems that needs shipping to a redis queue.

=== RabbitMQ

rabbitmqctl add_user logshipper logshipper
rabbitmqctl set_permissions logshipper "amqp-input" ".*" ".*"
