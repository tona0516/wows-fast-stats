<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Col, Container, FormGroup, Input, Label, Row } from "sveltestrap";
  import { OpenDirectory } from "wailsjs/go/main/App";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;

  const dispatch = createEventDispatcher();

  const openDirectory = async (path: string) => {
    try {
      await OpenDirectory(path);
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };
</script>

<Container fluid class="mt-2">
  <FormGroup>
    <Row>
      <Col
        ><Input
          type="switch"
          label="自動でスクリーンショットを保存する"
          bind:checked={inputUserConfig.save_screenshot}
          on:change={() => dispatch("Change")}
        /></Col
      >
    </Row>
    <Row>
      <Col>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click={() => openDirectory("screenshot/")}
          ><i class="bi bi-folder2-open" /> 保存フォルダを開く
        </a></Col
      >
    </Row>
  </FormGroup>
  <FormGroup>
    <Row>
      <Col>
        <Input
          type="switch"
          label="【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する"
          bind:checked={inputUserConfig.save_temp_arena_info}
          on:change={() => dispatch("Change")}
        /></Col
      >
    </Row>
    <Row>
      <Col>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click={() => openDirectory("temp_arena_info/")}
          ><i class="bi bi-folder2-open" /> 保存フォルダを開く
        </a></Col
      >
    </Row>
  </FormGroup>
  <FormGroup>
    <Row>
      <Col>
        <Input
          type="switch"
          label="アプリ改善のためのデータ送信を許可する"
          bind:checked={inputUserConfig.send_report}
          on:change={() => dispatch("Change")}
        /></Col
      >
    </Row>
    <Row>
      <Col>
        <ul>
          <li>アプリバージョン</li>
          <li>エラーログ</li>
          <li>設定値(config/user.json)</li>
          <li>戦闘情報(tempArenaInfo.json)</li>
        </ul></Col
      >
    </Row>
  </FormGroup>
  <FormGroup>
    <Row>
      <Col>
        <Input
          type="switch"
          label="新しいバージョンがある場合に通知する"
          bind:checked={inputUserConfig.notify_updatable}
          on:change={() => dispatch("Change")}
        /></Col
      >
    </Row>
  </FormGroup>
</Container>
