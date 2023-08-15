<script lang="ts">
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
  import { SkillLevelConverter } from "../RankConverter";
  import { Const } from "../Const";
  import type { domain } from "../../wailsjs/go/models";

  export let userConfig: domain.UserConfig;
</script>

<div class="center mt-2 mx-2">
  {#if userConfig.custom_color}
    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <th>スキル</th>
        <th colspan="2">PR</th>
        <th>ダメージ(平均比)</th>
        <th>勝率</th>
      </thead>
      <tbody>
        {#each Object.values(Const.SKILL_LEVELS) as v}
          <tr>
            <td>{v.level}</td>
            <td
              style="background-color: {SkillLevelConverter.fromPR(
                v.minPR,
                userConfig.custom_color.skill
              ).toBgColorCode()};">player_name</td
            >
            <td
              style="color: {SkillLevelConverter.fromPR(
                v.minPR,
                userConfig.custom_color.skill
              ).toTextColorCode()};">{v.minPR} ~ {v.maxPR}</td
            >
            <td
              style="color: {SkillLevelConverter.fromDamage(
                v.minDamage,
                1.0,
                userConfig.custom_color.skill
              ).toTextColorCode()};">{v.minDamage}倍 ~ {v.maxDamage}倍</td
            >
            <td
              style="color: {SkillLevelConverter.fromWinRate(
                v.minWin,
                userConfig.custom_color.skill
              ).toTextColorCode()};">{v.minWin}% ~ {v.maxWin}%</td
            >
          </tr>
        {/each}
      </tbody>
    </table>

    <ul>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click={() =>
            BrowserOpenURL("https://asia.wows-numbers.com/personal/rating")}
          >PRについて <i class="bi bi-box-arrow-up-right" /></a
        >
      </li>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click={() =>
            BrowserOpenURL(
              "https://asia.wows-numbers.com/personal/rating/expected/preview/"
            )}>艦種別平均値について <i class="bi bi-box-arrow-up-right" /></a
        >
      </li>
    </ul>
  {/if}
</div>
