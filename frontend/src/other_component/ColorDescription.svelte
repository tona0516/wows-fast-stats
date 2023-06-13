<script lang="ts">
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { SkillLevelConverter } from "../RankConverter";

const prColors: {
  label: string;
  minPR: number;
  maxPR: number;
  minDamage: number;
  maxDamage: number;
  minWin: number;
  maxWin: number;
}[] = [
  {
    label: "Bad",
    minPR: 0,
    maxPR: 750,
    minDamage: 0.0,
    maxDamage: 0.6,
    minWin: 0,
    maxWin: 47,
  },
  {
    label: "Below Average",
    minPR: 750,
    maxPR: 1100,
    minDamage: 0.6,
    maxDamage: 0.8,
    minWin: 47,
    maxWin: 50,
  },
  {
    label: "Average",
    minPR: 1100,
    maxPR: 1350,
    minDamage: 0.8,
    maxDamage: 1.0,
    minWin: 50,
    maxWin: 52,
  },
  {
    label: "Good",
    minPR: 1350,
    maxPR: 1550,
    minDamage: 1.0,
    maxDamage: 1.2,
    minWin: 52,
    maxWin: 54,
  },
  {
    label: "Very Good",
    minPR: 1550,
    maxPR: 1750,
    minDamage: 1.2,
    maxDamage: 1.4,
    minWin: 54,
    maxWin: 56,
  },
  {
    label: "Great",
    minPR: 1750,
    maxPR: 2100,
    minDamage: 1.4,
    maxDamage: 1.5,
    minWin: 56,
    maxWin: 60,
  },
  {
    label: "Unicum",
    minPR: 2100,
    maxPR: 2450,
    minDamage: 1.5,
    maxDamage: 1.6,
    minWin: 60,
    maxWin: 65,
  },
  {
    label: "Super Unicum",
    minPR: 2450,
    maxPR: 9999,
    minDamage: 1.6,
    maxDamage: 9999,
    minWin: 65,
    maxWin: 100,
  },
];
</script>

<h5 class="text-center">各項目の配色について</h5>

<table class="table table-sm table-text-color w-auto">
  <thead>
    <th>スキル</th>
    <th colspan="2">PR</th>
    <th>ダメージ(平均比)</th>
    <th>勝率</th>
  </thead>
  <tbody>
    {#each Object.values(prColors) as v}
      <tr>
        <td class="text-center">{v.label}</td>
        <td
          class="text-center"
          style="background-color: {SkillLevelConverter.fromPR(
            v.minPR
          ).toBgColorCode()};">player_name</td
        >
        <td
          class="text-center"
          style="color: {SkillLevelConverter.fromPR(
            v.minPR
          ).toTextColorCode()};">{v.minPR} ~ {v.maxPR}</td
        >
        <td
          class="text-center"
          style="color: {SkillLevelConverter.fromDamage(
            v.minDamage,
            1.0
          ).toTextColorCode()};">{v.minDamage}倍 ~ {v.maxDamage}倍</td
        >
        <td
          class="text-center"
          style="color: {SkillLevelConverter.fromWinRate(
            v.minWin
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
      on:click="{() =>
        BrowserOpenURL('https://asia.wows-numbers.com/personal/rating')}"
      >PRについて <i class="bi bi-box-arrow-up-right"></i></a
    >
  </li>
  <li>
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click="{() =>
        BrowserOpenURL(
          'https://asia.wows-numbers.com/personal/rating/expected/preview/'
        )}">艦種別平均値について <i class="bi bi-box-arrow-up-right"></i></a
    >
  </li>
</ul>
