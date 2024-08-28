package micro.evnt.test;

import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;

import java.io.IOException;
import java.util.List;

@Controller
@RequiredArgsConstructor
public class KafkaController {

    private static final Logger logger = LoggerFactory.getLogger(KafkaController.class);
    private final KafkaProducer kafkaProducer;
    private final ChessService chessService;

    @PostMapping("/send-email")
    @ResponseBody
    public String sendEmail(@RequestParam final String email) {
        logger.info("Received request to send email -> {}", email);
        kafkaProducer.sendMessage(email);
        logger.info("Email sent to email topic");
        kafkaProducer.saveToDbMessage(email);
        logger.info("Email sent to save to db");
        return "<p>Success: Email sent to 'email-topic' and 'savetodb' Kafka topics</p>";
    }

    @GetMapping("/")
    public String home(Model model) {
        try {
            List<Models.Streamer> topStreamers = chessService.getTopChessStreamersWithTwitchUrls();
            model.addAttribute("streamers", topStreamers);
        } catch (IOException e) {
            logger.error("Error fetching top streamers", e);
            return "error"; // Ensure there's an error page configured
        }
        return "index"; // Refers to the Thymeleaf template named 'index.html'
    }
}
