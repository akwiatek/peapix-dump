package io.github.akwiatek.peapixdump;

import static java.util.stream.Collectors.toList;
import static org.apache.commons.io.IOUtils.copy;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.net.URL;
import java.nio.file.Path;
import java.util.Date;
import java.util.List;

record FeedResponseItem(
    @JsonProperty String title,
    @JsonProperty String copyright,
    @JsonProperty URL fullUrl,
    @JsonProperty URL thumbUrl,
    @JsonProperty URL imageUrl,
    @JsonProperty Date date,
    @JsonProperty URL pageUrl
) {

  public URL getImageUrl() {
    return imageUrl;
  }
}

public class Main {

  public static final ObjectMapper MAPPER = new ObjectMapper();

  public static void main(String... args) throws Exception {
    var feedUrl = new URL("https://peapix.com/bing/feed");
    List<FeedResponseItem> items = MAPPER.readValue(feedUrl,
        new TypeReference<List<FeedResponseItem>>() {
        });
    items.parallelStream().peek(item -> download(item, items.indexOf(item), items.size()))
        .collect(toList());
  }

  private static void download(FeedResponseItem item, int index, int total) {
    var imageUrl = item.getImageUrl();
    var target = Path.of(imageUrl.getPath()).getFileName().toFile();
    try (var is = imageUrl.openStream(); var os = new FileOutputStream(target)) {
      copy(is, os);
      System.out.printf("%d/%d %s downloaded%n", index, total, imageUrl);
    } catch (IOException e) {
      throw new UncheckedIOException(e);
    }
  }
}
