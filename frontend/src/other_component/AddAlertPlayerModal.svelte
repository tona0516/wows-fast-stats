<script lang="ts">
  import Svelecte from "svelecte";
  import clone from "clone";
  import { createEventDispatcher } from "svelte";
  import {
    Modal,
    ModalBody,
    ModalFooter,
    Button,
    FormGroup,
    Label,
    Input,
  } from "sveltestrap";
  import {
    SearchPlayer,
    AlertPatterns,
    UpdateAlertPlayer,
  } from "../../wailsjs/go/main/App";
  import type { domain } from "../../wailsjs/go/models";
  import { Const } from "../Const";

  export const toggle = () => (open = !open);

  const dispatch = createEventDispatcher();

  let open = false;

  let alertPatterns: string[] = [];
  let target: domain.AlertPlayer = clone(Const.DEFAULT_ALERT_PLAYER);
  let searchResult: domain.WGAccountListData;

  async function search(input: string) {
    const accountList = await SearchPlayer(input);
    return accountList.data;
  }

  async function onOpen() {
    alertPatterns = await AlertPatterns();
    target = clone(Const.DEFAULT_ALERT_PLAYER);
    searchResult = undefined;
  }

  async function add(
    player: domain.AlertPlayer,
    searchResult: domain.WGAccountListData
  ) {
    try {
      if (!searchResult) {
        dispatch("Failure", { message: "不正な入力です" });
        return;
      }
      player.account_id = searchResult.account_id;
      player.name = searchResult.nickname;

      if (!validate(player)) {
        dispatch("Failure", { message: "不正な入力です" });
        return;
      }

      await UpdateAlertPlayer(player);
      dispatch("Success");
    } catch (error) {
      dispatch("Failure", { message: error });
      return;
    } finally {
      toggle();
    }
  }

  function validate(player: domain.AlertPlayer): boolean {
    return (
      player.account_id !== 0 && player.name !== "" && player.pattern !== ""
    );
  }
</script>

<Modal isOpen={open} {toggle} on:open={onOpen}>
  <ModalBody class="modal-color">
    <FormGroup>
      <Label>プレイヤー名</Label>
      <Svelecte
        style="color: #2d2c2c;"
        valueAsObject
        id="player"
        fetch={search}
        placeholder=""
        minQuery="3"
        labelField="nickname"
        bind:value={searchResult}
      />
    </FormGroup>

    <FormGroup>
      <Label>アイコン</Label>
      <!-- TODO migrate sveltestrap -->
      <div>
        {#each alertPatterns as pattern}
          <div class="form-check form-check-inline">
            <input
              class="form-check-input"
              type="radio"
              bind:group={target.pattern}
              value={pattern}
            />
            <label class="form-check-label" for="icon">
              <i class="bi {pattern}" />
            </label>
          </div>
        {/each}
      </div>
    </FormGroup>

    <FormGroup>
      <Label>メモ(任意)</Label>
      <Input
        type="textarea"
        maxlength={Const.MAX_MEMO_LENGTH}
        bind:value={target.message}
      />
      <span>{target.message.length}/{Const.MAX_MEMO_LENGTH}</span>
    </FormGroup>
  </ModalBody>

  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click={toggle}>キャンセル</Button>
    <Button size="sm" color="primary" on:click={() => add(target, searchResult)}
      >追加</Button
    >
  </ModalFooter>
</Modal>
