<script lang="ts">
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    Alert,
    Badge,
    Button,
    FormGroup,
    Input,
    Label,
    Spinner,
  } from "sveltestrap";
  import {
    ApplyRequiredUserConfig,
    SelectDirectory,
    UserConfig,
  } from "wailsjs/go/main/App";
  import { domain, vo } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;

  let isLoading = false;
  let validatedResult: vo.ValidatedResult;

  const dispatch = createEventDispatcher();

  const selectDirectory = async () => {
    try {
      const path = await SelectDirectory();
      if (!path) return;
      inputUserConfig.install_path = path;
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };

  const clickApply = async () => {
    isLoading = true;
    try {
      validatedResult = await ApplyRequiredUserConfig(
        inputUserConfig.install_path,
        inputUserConfig.appid,
      );

      const errorTexts = Object.values(validatedResult);
      const isValid =
        errorTexts.filter((it) => it == "").length === errorTexts.length;
      if (isValid) {
        const latest = await UserConfig();
        storedUserConfig.set(latest);
        dispatch("UpdateSuccess", { message: "設定を更新しました。" });
      }
    } catch (error) {
      inputUserConfig = $storedUserConfig;
      dispatch("Failure", { message: error });
    } finally {
      isLoading = false;
    }
  };
</script>

<FormGroup>
  <Label
    >World of Warshipsインストールフォルダ <Badge color="danger">必須</Badge>
  </Label>
  <div class="d-flex justify-content-center">
    <Input type="text" bind:value={inputUserConfig.install_path} />
    <Button color="secondary" on:click={selectDirectory}>選択</Button>
  </div>

  {#if validatedResult?.install_path}
    <Alert color="danger" class="m-1">
      {validatedResult?.install_path}
    </Alert>
  {/if}
</FormGroup>

<FormGroup>
  <Label>アプリケーションID <Badge color="danger">必須</Badge></Label>
  <Input type="text" bind:value={inputUserConfig.appid} />
  <div>
    <ExternalLink
      url="https://developers.wargaming.net/"
      text="Developer Room"
    />で作成したIDを入力してください。
  </div>

  {#if validatedResult?.appid}
    <Alert color="danger" class="m-1">
      {validatedResult?.appid}
    </Alert>
  {/if}
</FormGroup>

<!-- apply -->
<FormGroup class="center">
  <Button size="sm" color="primary" disabled={isLoading} on:click={clickApply}>
    {#if isLoading}
      <Spinner size="sm" />更新中...
    {:else}
      保存
    {/if}
  </Button>
</FormGroup>
