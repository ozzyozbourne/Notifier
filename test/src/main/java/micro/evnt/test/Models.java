package micro.evnt.test;

import lombok.Data;
import lombok.NoArgsConstructor;

@NoArgsConstructor
abstract public class Models {


    @Data
    public static class Player {
        private String username;
    }

    @Data
    public static class Streamer {
        private String username;
        private String twitch_url;
        private String avatar;
    }

}
