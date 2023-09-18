<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { FormGroup, Label, Input, Col, Container, Row } from "sveltestrap";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;
  export let defaultUserConfig: domain.UserConfig;

  const dispatch = createEventDispatcher();
</script>

<Container fluid class="mt-2">
  <FormGroup>
    <Row>
      <Col><Label>チーム平均に含める最小戦闘数(艦戦績)</Label></Col>
    </Row>
    <Row>
      <Col>
        <Input
          type="range"
          min={defaultUserConfig.team_average.min_ship_battles}
          max={100}
          step={1}
          bind:value={inputUserConfig.team_average.min_ship_battles}
          on:change={() => dispatch("Change")}
        />
        {inputUserConfig.team_average.min_ship_battles}</Col
      >
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col><Label>チーム平均に含める最小戦闘数(総合成績)</Label></Col>
    </Row>
    <Row>
      <Col>
        <Input
          type="range"
          min={defaultUserConfig.team_average.min_overall_battles}
          max={3000}
          step={10}
          bind:value={inputUserConfig.team_average.min_overall_battles}
          on:change={() => dispatch("Change")}
        />
        {inputUserConfig.team_average.min_overall_battles}</Col
      >
    </Row>
  </FormGroup>
</Container>
