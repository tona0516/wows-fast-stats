<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { FormGroup, Label, Input } from "sveltestrap";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;
  export let defaultUserConfig: domain.UserConfig;

  const dispatch = createEventDispatcher();
</script>

<FormGroup>
  <Label>チーム平均に含める最小戦闘数(艦戦績)</Label>
  <Input
    type="range"
    min={defaultUserConfig.team_average.min_ship_battles}
    max={100}
    step={1}
    bind:value={inputUserConfig.team_average.min_ship_battles}
    on:change={() => dispatch("Change")}
  />
  {inputUserConfig.team_average.min_ship_battles}
</FormGroup>

<FormGroup>
  <Label>チーム平均に含める最小戦闘数(総合成績)</Label>
  <Input
    type="range"
    min={defaultUserConfig.team_average.min_overall_battles}
    max={3000}
    step={10}
    bind:value={inputUserConfig.team_average.min_overall_battles}
    on:change={() => dispatch("Change")}
  />
  {inputUserConfig.team_average.min_overall_battles}
</FormGroup>
