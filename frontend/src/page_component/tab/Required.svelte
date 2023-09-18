<script lang="ts">
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    Button,
    Col,
    Container,
    FormGroup,
    Input,
    Label,
    Row,
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

<Container fluid class="mt-2">
  <FormGroup>
    <Row>
      <Col><Label>World of Warshipsインストールフォルダ</Label></Col>
    </Row>
    <Row>
      <Col>
        <div class="d-flex justify-content-center">
          <Input
            bind:value={inputUserConfig.install_path}
            feedback={validatedResult?.install_path}
            invalid={validatedResult !== undefined &&
              validatedResult.install_path !== ""}
          />
          <Button color="secondary" on:click={selectDirectory}>選択</Button>
        </div></Col
      >
    </Row>
    <Row>
      <Col
        ><i class="bi bi-info-circle-fill" /> ゲームクライアントの実行ファイルがあるフォルダを選択してください。</Col
      >
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col><Label>アプリケーションID</Label></Col>
    </Row>
    <Row>
      <Col
        ><Input
          bind:value={inputUserConfig.appid}
          feedback={validatedResult?.appid}
          invalid={validatedResult !== undefined &&
            validatedResult.appid !== ""}
        /></Col
      >
    </Row>
    <Row>
      <Col
        ><i class="bi bi-info-circle-fill" />
        <ExternalLink
          url="https://developers.wargaming.net/"
          text="Developer Room"
        />で作成したIDを入力してください。</Col
      >
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col sm={{ size: 2, offset: 5 }}>
        <Button color="primary" disabled={isLoading} on:click={clickApply}>
          {#if isLoading}
            <Spinner size="sm" /> 更新中...
          {:else}
            保存
          {/if}
        </Button>
      </Col>
    </Row>
  </FormGroup>
</Container>
