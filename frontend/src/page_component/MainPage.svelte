<script lang="ts">
  import { MAIN_PAGE_ID } from "src/lib/types";
  import BattleMeta from "src/other_component/BattleMeta.svelte";
  import ColorDescription from "src/other_component/ColorDescription.svelte";
  import Ofuse from "src/other_component/Ofuse.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
  import SummaryGraph from "src/other_component/SummaryGraph.svelte";
  import { storedBattle, storedUserConfig, storedSummary } from "src/stores";
  import { Col, Container, Row } from "sveltestrap";
</script>

<Container fluid id={MAIN_PAGE_ID} class="mt-2">
  <Row>
    <Col>
      {#if $storedBattle?.meta}
        <BattleMeta meta={$storedBattle?.meta} />
      {/if}
    </Col>
  </Row>
  <Row>
    <Col>
      {#if $storedBattle}
        <StatisticsTable
          teams={$storedBattle.teams}
          userConfig={$storedUserConfig}
          on:EditAlertPlayer
          on:RemoveAlertPlayer
          on:CheckPlayer
        />
      {/if}
    </Col>
  </Row>
  <Row>
    <Col sm={{ size: 6, offset: 3 }}>
      {#if $storedSummary}
        <SummaryGraph summary={$storedSummary} />
      {/if}
    </Col>
  </Row>
  <Row>
    <Col>
      <ColorDescription userConfig={$storedUserConfig} />
    </Col>
  </Row>
  <Row>
    <Col>
      <Ofuse />
    </Col>
  </Row>
</Container>
