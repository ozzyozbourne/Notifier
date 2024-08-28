package micro.evnt.test;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Service
public class ChessService {

    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;

    @Autowired
    public ChessService(RestTemplate restTemplate, ObjectMapper objectMapper) {
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
    }

    public List<Models.Streamer> getTopChessStreamersWithTwitchUrls() throws IOException {
        List<Models.Player> dailyTopPlayers = fetchDailyTopPlayers();
        List<Models.Streamer> streamers = fetchStreamerData();

        // Debug logs to verify data
        System.out.println("Daily Top Players: " + dailyTopPlayers);
        System.out.println("Streamers: " + streamers);

        // Find players who are in both the daily top and streamers with Twitch URLs
        List<Models.Streamer> topStreamers = streamers.stream()
                .filter(streamer -> dailyTopPlayers.stream()
                        .anyMatch(player -> player.getUsername().equalsIgnoreCase(streamer.getUsername()))
                        && streamer.getTwitch_url() != null && !streamer.getUsername().isEmpty())
                .collect(Collectors.toList());

        System.out.println("Top Streamers with Twitch URLs: " + topStreamers);

        return topStreamers;
    }

    private List<Models.Player> fetchDailyTopPlayers() throws IOException {
        String url = "https://api.chess.com/pub/leaderboards";
        String jsonResponse = restTemplate.getForObject(url, String.class);

        Map<String, List<Models.Player>> response = objectMapper.readValue(
                jsonResponse, new TypeReference<>() {
                });

        return response.getOrDefault("daily", List.of());
    }

    private List<Models.Streamer> fetchStreamerData() throws IOException {
        String url = "https://api.chess.com/pub/streamers";
        String jsonResponse = restTemplate.getForObject(url, String.class);

        Map<String, List<Models.Streamer>> response = objectMapper.readValue(
                jsonResponse, new TypeReference<>() {
                });

        return response.getOrDefault("streamers", List.of());
    }
}
