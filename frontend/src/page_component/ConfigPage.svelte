<script lang="ts">
  import clone from "clone";
  import { get } from "svelte/store";
  import { storedUserConfig } from "src/stores";
  import { UserConfig, ApplyUserConfig } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import Required from "./tab/Required.svelte";
  import Other from "./tab/Other.svelte";
  import TeamSummary from "./tab/TeamSummary.svelte";
  import Display from "./tab/Display.svelte";
  import { TabContent, TabPane } from "sveltestrap";

  export let defaultUserConfig: domain.UserConfig;

  let inputUserConfig = clone($storedUserConfig);
  storedUserConfig.subscribe((it) => {
    inputUserConfig = clone(it);
  });

  const silentApply = async () => {
    // Note: for the following sveltestrap bug
    // https://github.com/bestguy/sveltestrap/issues/461
    await new Promise((resolve) => setTimeout(resolve, 100));

    try {
      await ApplyUserConfig(inputUserConfig);
      const latest = await UserConfig();
      storedUserConfig.set(latest);
    } catch (error) {
      inputUserConfig = get(storedUserConfig);
    }
  };
</script>

<TabContent>
  <TabPane tabId="required" tab="必須設定" active>
    <Required {inputUserConfig} on:UpdateSuccess on:Failure />
  </TabPane>
  <TabPane tabId="display" tab="表示設定">
    <Display {inputUserConfig} {defaultUserConfig} on:Change={silentApply} />
  </TabPane>
  <TabPane tabId="team-summary" tab="チームサマリー設定">
    <TeamSummary
      {inputUserConfig}
      {defaultUserConfig}
      on:Change={silentApply}
    />
  </TabPane>
  <TabPane tabId="other" tab="その他設定">
    <Other {inputUserConfig} on:Change={silentApply} on:Failure />
  </TabPane>
</TabContent>
