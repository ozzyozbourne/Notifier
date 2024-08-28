package micro.evnt.test;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestTemplate;

@Configuration
public class KafkaTopicConfig {

    private final String emailTopicName;
    private final String saveToDbTopicName;

    @Autowired
    public KafkaTopicConfig(final @Value("${spring.kafka.topic.email-topic}") String emailTopicName,
                            final @Value("${spring.kafka.topic.savetodb}") String saveToDbTopicName) {
        this.emailTopicName = emailTopicName;
        this.saveToDbTopicName = saveToDbTopicName;
    }


    @Bean
    public NewTopic emailTopic() {
        return new NewTopic(emailTopicName, 1, (short) 1);
    }

    @Bean
    public NewTopic saveToDbTopic() {
        return new NewTopic(saveToDbTopicName, 1, (short) 1);
    }

    @Bean
    public RestTemplate restTemplate() {
        return new RestTemplate();
    }
}
