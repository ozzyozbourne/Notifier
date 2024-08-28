package micro.evnt.test;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.annotation.CreatedDate;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.AllArgsConstructor;
import java.time.Instant;

@Document(collection = "emails")
@Data
@NoArgsConstructor
@AllArgsConstructor
public class EmailDocument {

    @Id
    private String id;

    private String email;

    @CreatedDate
    private Instant createdAt;
}
