<script lang="ts">
  import { data } from "wailsjs/go/models";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { DispName } from "src/lib/DispName";
  import { RATING_DEFS, RatingInfo } from "src/lib/RatingLevel";
  import { THREAT_LEVEL_DEFS } from "src/lib/ThreatLevel";

  export let config: data.UserConfigV2;

  interface RatingDescription {
    ratingText: string;
    playerName: string;
    playerNameColor: string;
    prRange: string;
    prColor: string;
    damageRange: string;
    damageColor: string;
    winRateRange: string;
    winRateColor: string;
  }

  interface ThreatDescription {
    threatText: string;
    range: string;
    textColor: string;
    bgColor: string;
  }

  const COLUMNS = [
    { text: "スキル", colspan: 1 },
    { text: "PR", colspan: 2 },
    { text: "ダメージ(平均比)", colspan: 1 },
    { text: "勝率", colspan: 1 },
  ];

  const COLUMNS2 = [
    { text: "スキル", colspan: 1 },
    { text: "戦力評価", colspan: 1 },
  ];

  const LINKS = [
    {
      url: "https://asia.wows-numbers.com/personal/rating",
      text: "算出方法",
    },
    {
      url: "https://asia.wows-numbers.com/personal/rating/expected/preview/",
      text: "艦種別平均値について",
    },
  ];

  const getRatingDesc = (
    skillColorCode: data.UCSkillColorCode,
  ): RatingDescription[] => {
    const result: RatingDescription[] = [];

    for (let i = 0; i < RATING_DEFS.length; i++) {
      const current = RATING_DEFS[i];

      let prRange, damageRange, winRateRange;
      if (RATING_DEFS[i + 1]) {
        const next = RATING_DEFS[i + 1];
        prRange = `${current.pr} ~ ${next.pr}`;
        damageRange = `${current.damage} ~ ${next.damage}倍`;
        winRateRange = `${current.winRate} ~ ${next.winRate}%`;
      } else {
        prRange = `${current.pr} ~`;
        damageRange = `${current.damage}倍 ~`;
        winRateRange = `${current.winRate}% ~`;
      }

      result.push({
        ratingText: DispName.SKILL_LEVELS.get(current.level)!,
        playerName: "player_name",
        playerNameColor: RatingInfo.fromPR(current.pr, skillColorCode)!
          .textColorCode,
        prRange: prRange,
        prColor: RatingInfo.fromDamage(current.damage, 1.0, skillColorCode)!
          .textColorCode,
        damageRange: damageRange,
        damageColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)!
          .textColorCode,
        winRateRange: winRateRange,
        winRateColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)!
          .textColorCode,
      });
    }

    return result;
  };

  const getThreatDesc = (): ThreatDescription[] => {
    const result: ThreatDescription[] = [];

    for (let i = 0; i < THREAT_LEVEL_DEFS.length; i++) {
      const current = THREAT_LEVEL_DEFS[i];

      let range;
      if (THREAT_LEVEL_DEFS[i + 1]) {
        const next = THREAT_LEVEL_DEFS[i + 1];
        range = `${current.score} ~ ${next.score}`;
      } else {
        range = `${current.score} ~`;
      }

      result.push({
        threatText: current.info.level,
        range: range,
        textColor: current.info.textColorCode,
        bgColor: current.info.bgColorCode,
      });
    }

    return result;
  };

  const ratingDesc = getRatingDesc(config.color.skill.text);
  const threatDesc = getThreatDesc();
</script>

{#if config.color}
  <h5>Personal Rating (by WoWS Numbers)</h5>

  <UkTable>
    <thead>
      {#each COLUMNS as column}
        <th class="uk-text-center" colspan={column.colspan}>{column.text}</th>
      {/each}
    </thead>

    <tbody>
      {#each ratingDesc as desc}
        <tr>
          <td class="uk-text-center">{desc.ratingText}</td>
          <td class="uk-text-center" style="color: {desc.playerNameColor};"
            >{desc.playerName}</td
          >
          <td class="uk-text-center" style="color: {desc.prColor};"
            >{desc.prRange}</td
          >
          <td class="uk-text-center" style="color: {desc.damageColor};"
            >{desc.damageRange}</td
          >
          <td class="uk-text-center" style="color: {desc.winRateColor};"
            >{desc.winRateRange}</td
          >
        </tr>
      {/each}
    </tbody>
  </UkTable>

  {#each LINKS as link}
    <div class="uk-margin-small-right">
      <ExternalLink url={link.url}>{link.text}</ExternalLink>
    </div>
  {/each}
{/if}

<h5>戦力評価 (by 178usagi氏)</h5>

<UkTable>
  <thead>
    {#each COLUMNS2 as column}
      <th class="uk-text-center" colspan={column.colspan}>{column.text}</th>
    {/each}
  </thead>

  <tbody>
    {#each threatDesc as desc}
      <tr>
        <td class="uk-text-center">{desc.threatText}</td>
        <td
          class="uk-text-center"
          style="color: {desc.textColor}; background-color: {desc.bgColor}"
          >{desc.range}</td
        >
      </tr>
    {/each}
  </tbody>
</UkTable>
