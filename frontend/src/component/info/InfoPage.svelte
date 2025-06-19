<script lang="ts">
import iconApp from "src/assets/images/appicon.png";
import ExternalLink from "src/component/common/ExternalLink.svelte";
import type {
  DispColorTableInfo,
  Row,
} from "src/component/info/internal/DispColorTableInfo";
import { DispName } from "src/lib/DispName";
import { RATING_DEFS, type RatingDef, RatingInfo } from "src/lib/RatingLevel";
import { THREAT_LEVEL_DEFS, type ThreatLevelDef } from "src/lib/ThreatLevel";
import { storedConfig } from "src/stores";
import { Semver } from "wailsjs/go/main/App";
import type { data } from "wailsjs/go/models";
import DispColorTable from "./internal/DispColorTable.svelte";

const LINKS = [
  {
    icon: "question-circle",
    url: "https://github.com/tona0516/wows-fast-stats/wiki/FAQ",
    text: "FAQ",
  },
  {
    icon: "twitter",
    url: "https://twitter.com/tonango_0516",
    text: "@tonango_0516",
  },
  {
    icon: "github",
    url: "https://github.com/tona0516/wows-fast-stats",
    text: "tona0516/wows-fast-stats",
  },
];

const getPRTableInfo = (
  defs: RatingDef[],
  skillColorCode: data.UCSkillColorCode,
): DispColorTableInfo => {
  const rows: Row[] = [];

  for (let i = 0; i < defs.length; i++) {
    const current = defs[i];

    let prRange = "";
    let damageRange = "";
    let winRateRange = "";

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
        text: DispName.SKILL_LEVELS.get(current.level) ?? "",
      },
      {
        text: prRange,
        textColor: RatingInfo.fromDamage(current.damage, 1.0, skillColorCode)
          ?.textColorCode,
      },
      {
        text: damageRange,
        textColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)
          ?.textColorCode,
      },
      {
        text: winRateRange,
        textColor: RatingInfo.fromWinRate(current.winRate, skillColorCode)
          ?.textColorCode,
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

    let range = "";
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
      },
    ];

    rows.push(row);
  }

  return { headers: ["スキル", "戦力評価"], rows };
};

const prTableInfo = getPRTableInfo(RATING_DEFS, $storedConfig.color.skill.text);
const tlTableInfo = getTLTableInfo(THREAT_LEVEL_DEFS);
</script>

<div>
  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">Personal Rating (by WoWS Numbers)</p>
    <DispColorTable tableInfo={prTableInfo} />
    <ExternalLink url={"https://asia.wows-numbers.com/personal/rating"}
      >Personal Ratingの算出方法</ExternalLink
    >
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">戦力評価 (by 178usagi氏))</p>
    <DispColorTable tableInfo={tlTableInfo} />
    <ExternalLink url={"https://github.com/tona0516/wows-fast-stats/wiki"}
      >戦力評価の算出方法</ExternalLink
    >
  </div>

  <div class="p-4 flex flex-col items-center">
    <img src={iconApp} alt="" width="128px" height="128px" />
    <div class="pt-1">
      wows-fast-stats {#await Semver() then semver} {semver} {/await}
    </div>
    <div class="pt-2">
      {#each LINKS as link}
        <div class="flex flex-col items-center">
          <ExternalLink url={link.url}>
            <i class="bi bi-{link.icon}">{link.text}</i>
          </ExternalLink>
        </div>
      {/each}
    </div>
    <div class="pt-2">Copyright © 2023 tona0516 All Rights Reserved.</div>
  </div>
</div>
