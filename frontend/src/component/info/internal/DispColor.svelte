<script lang="ts">
  import { data } from "wailsjs/go/models";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import { DispName } from "src/lib/DispName";
  import { RATING_DEFS, RatingInfo, type RatingDef } from "src/lib/RatingLevel";
  import { THREAT_LEVEL_DEFS, type ThreatLevelDef } from "src/lib/ThreatLevel";
  import type {
    DispColorTableInfo,
    Row,
  } from "src/component/info/internal/DispColorTableInfo";
  import DispColorTable from "./DispColorTable.svelte";

  export let config: data.UserConfigV2;

  const LINKS = [
    {
      url: "https://asia.wows-numbers.com/personal/rating",
      text: "Personal Ratingの算出方法",
    },
    {
      url: "https://github.com/tona0516/wows-fast-stats/wiki",
      text: "戦力評価の算出方法",
    },
    {
      url: "https://asia.wows-numbers.com/personal/rating/expected/preview/",
      text: "艦種別平均値について",
    },
  ];

  const getPRTableInfo = (
    defs: RatingDef[],
    skillColorCode: data.UCSkillColorCode,
  ): DispColorTableInfo => {
    const rows: Row[] = [];

    for (let i = 0; i < defs.length; i++) {
      const current = defs[i];

      let prRange, damageRange, winRateRange;
      if (defs[i + 1]) {
        const next = defs[i + 1];
        prRange = `${current.pr} ~ ${next.pr}`;
        damageRange = `${current.damage} ~ ${next.damage}倍`;
        winRateRange = `${current.winRate} ~ ${next.winRate}%`;
      } else {
        prRange = `${current.pr} ~`;
        damageRange = `${current.damage}倍 ~`;
        winRateRange = `${current.winRate}% ~`;
      }

      const row: Row = [
        {
          text: DispName.SKILL_LEVELS.get(current.level)!,
        },
        {
          text: prRange,
          textColor: RatingInfo.fromDamage(current.damage, 1.0, skillColorCode)!
            .textColorCode,
        },
        {
          text: damageRange,
          textColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)!
            .textColorCode,
        },
        {
          text: winRateRange,
          textColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)!
            .textColorCode,
        },
      ];

      rows.push(row);
    }

    return { headers: ["スキル", "PR", "ダメージ(平均比)", "勝率"], rows };
  };

  const getTLTableInfo = (defs: ThreatLevelDef[]): DispColorTableInfo => {
    const rows: Row[] = [];
    for (let i = 0; i < defs.length; i++) {
      const current = defs[i];

      let range;
      if (defs[i + 1]) {
        const next = defs[i + 1];
        range = `${current.score} ~ ${next.score}`;
      } else {
        range = `${current.score} ~`;
      }

      const row: Row = [
        {
          text: current.info.level,
        },
        {
          text: range,
          textColor: current.info.textColorCode,
          bgColor: current.info.bgColorCode,
        },
      ];

      rows.push(row);
    }

    return { headers: ["スキル", "戦力評価"], rows };
  };

  const prTableInfo = getPRTableInfo(RATING_DEFS, config.color.skill.text);
  const tlTableInfo = getTLTableInfo(THREAT_LEVEL_DEFS);
</script>

<div class="uk-padding-small">
  <DispColorTable
    caption="Personal Rating (by WoWS Numbers)"
    tableInfo={prTableInfo}
  />
</div>

<div class="uk-padding-small">
  <DispColorTable caption="戦力評価 (by 178usagi氏)" tableInfo={tlTableInfo} />
</div>

<div class="uk-padding-small">
  <ul>
    {#each LINKS as link}
      <li>
        <ExternalLink url={link.url}>{link.text}</ExternalLink>
      </li>
    {/each}
  </ul>
</div>
