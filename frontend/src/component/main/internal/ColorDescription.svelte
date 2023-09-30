<script lang="ts">
  import { colorDescriptions } from "src/lib/rating/RatingConst";
  import type { domain } from "wailsjs/go/models";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkTable from "src/component/common/uikit/UkTable.svelte";

  export let userConfig: domain.UserConfig;

  const COLUMNS = [
    { text: "スキル", colspan: 1 },
    { text: "PR", colspan: 2 },
    { text: "ダメージ(平均比)", colspan: 1 },
    { text: "勝率", colspan: 1 },
  ];

  const LINKS = [
    {
      url: "https://asia.wows-numbers.com/personal/rating",
      text: "PRについて",
    },
    {
      url: "https://asia.wows-numbers.com/personal/rating/expected/preview/",
      text: "艦種別平均値について",
    },
  ];
</script>

{#if userConfig.custom_color}
  <div class="uk-flex uk-flex-center">
    <UkTable>
      <thead>
        {#each COLUMNS as column}
          <th class="uk-text-center" colspan={column.colspan}>{column.text}</th>
        {/each}
      </thead>

      <tbody>
        {#each colorDescriptions(userConfig) as desc}
          <tr>
            <td class="uk-text-center">{desc.level.text}</td>
            <td
              class="uk-text-center"
              style="background-color: {desc.playerName.bgColor};"
              >{desc.playerName.text}</td
            >
            <td class="uk-text-center" style="color: {desc.pr.textColor};"
              >{desc.pr.text}</td
            >
            <td class="uk-text-center" style="color: {desc.damage.textColor};"
              >{desc.damage.text}</td
            >
            <td class="uk-text-center" style="color: {desc.win.textColor};"
              >{desc.win.text}</td
            >
          </tr>
        {/each}
      </tbody>
    </UkTable>
  </div>

  <div class="uk-flex uk-flex-center">
    {#each LINKS as link}
      <div class="uk-margin-small-right">
        <ExternalLink url={link.url}>{link.text}</ExternalLink>
      </div>
    {/each}
  </div>
{/if}
