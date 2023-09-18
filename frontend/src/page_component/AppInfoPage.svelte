<script lang="ts">
  import { Semver } from "wailsjs/go/main/App";
  import iconApp from "src/assets/images/appicon.png";
  import Logging from "src/other_component/Logging.svelte";
  import { storedLogs } from "src/stores";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import { Container, Row, Col, Icon } from "sveltestrap";
</script>

<Container fluid class="mt-2">
  <Row>
    <Col class="center my-1">
      <img src={iconApp} alt="" width="128px" height="128px" />
      <div class="d-flex justify-content-center">
        wows-fast-stats
        {#await Semver() then semver} v{semver} {/await}
      </div>
    </Col>
  </Row>
  <Row>
    <Col class="center my-1">
      <div>
        <Icon name="question-circle" />
        <ExternalLink
          url="https://github.com/tona0516/wows-fast-stats/wiki/FAQ"
          text="FAQ"
        />
      </div>
      <div>
        <Icon name="twitter" />
        <ExternalLink
          url="https://twitter.com/tonango_0516"
          text="@tonango_0516"
        />
      </div>
      <div>
        <Icon name="github" />
        <ExternalLink
          url="https://github.com/tona0516/wows-fast-stats"
          text="tona0516/wows-fast-stats"
        />
      </div>
      <div>Copyright © 2023 tona0516 All Rights Reserved.</div>
    </Col>
  </Row>
  <Row>
    <Col class="my-1">
      <Logging logs={$storedLogs} />
    </Col>
  </Row>
</Container>
