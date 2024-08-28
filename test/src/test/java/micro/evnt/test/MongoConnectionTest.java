package micro.evnt.test;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.data.mongo.DataMongoTest;
import org.springframework.test.context.ContextConfiguration;

import static org.assertj.core.api.Assertions.assertThat;

@DataMongoTest
@ContextConfiguration(classes = TestApplication.class)
public class MongoConnectionTest {

    @Autowired
    private EmailRepository emailRepository;

    @Test
    public void testMongoConnection() {
        EmailDocument emailDocument = new EmailDocument();
        emailDocument.setEmail("test@example.com");

        EmailDocument savedDocument = emailRepository.save(emailDocument);

        EmailDocument retrievedDocument = emailRepository.findById(savedDocument.getId()).orElse(null);

        assertThat(retrievedDocument).isNotNull();
        assertThat(retrievedDocument.getEmail()).isEqualTo("test@example.com");
    }
}
