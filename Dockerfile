FROM rabbitmq:3.13-management

# Enable STOMP plugin (just for demonstration)
RUN rabbitmq-plugins enable rabbitmq_stomp