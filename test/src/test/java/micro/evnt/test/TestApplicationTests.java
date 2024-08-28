package micro.evnt.test;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.times;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.kafka.test.context.EmbeddedKafka;

import org.springframework.mail.SimpleMailMessage;

import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
@EmbeddedKafka(partitions = 1, topics = {"email-topic"})
public class TestApplicationTests {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private KafkaProducer kafkaProducer;

    @MockBean
    private JavaMailSender mailSender;

    @Test
    public void testSendEmail() throws Exception {
        String email = "user@example.com";
        mockMvc.perform(MockMvcRequestBuilders.post("/api/v1/kafka/send-email")
                        .param("email", email))
                .andExpect(status().isOk());
        verify(kafkaProducer, times(1)).sendMessage(email);
        kafkaProducer.sendMessage(email);
        verify(mailSender, times(1)).send(any(SimpleMailMessage.class));
    }
}
