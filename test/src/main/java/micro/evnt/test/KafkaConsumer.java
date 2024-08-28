package micro.evnt.test;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

@Service
public final class KafkaConsumer {

    private static final Logger logger = LoggerFactory.getLogger(KafkaConsumer.class);

    private final JavaMailSender mailSender;
    private final String emailUsername;
    private final SimpleMailMessage message;
    private final EmailRepository emailRepository;


    @Autowired
    public KafkaConsumer(final JavaMailSender mailSender,
                         final @Value("${spring.mail.username}") String emailUsername,
                         final EmailRepository emailRepository)
    {
        this.mailSender = mailSender;
        this.emailUsername = emailUsername;
        this.emailRepository = emailRepository;
        this.message = new SimpleMailMessage();
    }

    @KafkaListener(topics = "email-topic", groupId = "email-group")
    public void listen(final String emailUsername) {
        logger.info("Received email: {} for email topic ", emailUsername);
        sendEmail(emailUsername, "Welcome!", "This is a test email sent through Kafka listener.");
    }

    @KafkaListener(topics = "savetodb", groupId = "email-group")
    public void listenToSaveToDbTopic(final String email) {
        logger.info("Received email for saving: {} for db", email);
        //saveToDatabase(email);
    }


    public void sendEmail(final String to, final String subject, final String body) {
        message.setTo(to);
        message.setSubject(subject);
        message.setText(body);
        message.setFrom(emailUsername);

        try {
            mailSender.send(message);
            logger.info("Email sent successfully to {}", to);
        } catch (final Exception e) {
            logger.error("Failed to send email: ->", e);
        }
    }

    private void saveToDatabase(final String email) {
        EmailDocument emailDocument = new EmailDocument();
        emailDocument.setEmail(email);
        logger.info("Saving to database: {}", email);
        try {
            emailRepository.save(emailDocument);
            logger.info("Email saved to database: {}", email);
        } catch (final Exception e) {
            logger.error("Failed to save email to database: ->", e);
        }
    }
}
