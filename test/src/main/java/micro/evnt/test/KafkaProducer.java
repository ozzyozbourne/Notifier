package micro.evnt.test;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public final class KafkaProducer {

    private static final Logger logger = LoggerFactory.getLogger(KafkaProducer.class);

    private final KafkaTemplate<String, String> kafkaTemplate;

    public void sendMessage(final String message) {
        logger.info("Sending message: {}", message);
        kafkaTemplate.send("email-topic", message);
        logger.info("Sent message: {}", message);
    }

    public void saveToDbMessage(final String message) {
        logger.info("Sending save message to 'savetodb' topic: {}", message);
        kafkaTemplate.send("savetodb", message);
        logger.info("Sent save message: {}", message);
    }
}