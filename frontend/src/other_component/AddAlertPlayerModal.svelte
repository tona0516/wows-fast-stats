<script lang="ts">
  // @ts-ignore
  import Svelecte from "svelecte";
  import clone from "clone";
  import { createEventDispatcher } from "svelte";
  import {
    Modal,
    ModalBody,
    FormGroup,
    Label,
    Input,
    ModalFooter,
    Button,
  } from "sveltestrap";
  import {
    SearchPlayer,
    AlertPatterns,
    UpdateAlertPlayer,
  } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";
  import { MAX_MEMO_LENGTH } from "src/const";

  export const toggle = () => (open = !open);

  let open = false;
  let target: domain.AlertPlayer;
  let searchResult: domain.WGAccountListData | undefined;

  const DEFAULT_ALERT_PLAYER: domain.AlertPlayer = {
    account_id: 0,
    name: "",
    pattern: "bi-check-circle-fill",
    message: "",
  };

  const dispatch = createEventDispatcher();

  const onOpen = async () => {
    target = clone(DEFAULT_ALERT_PLAYER);
    searchResult = undefined;
  };

  const add = async (player: domain.AlertPlayer) => {
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
  };

  const validate = (player: domain.AlertPlayer): boolean => {
    return (
      player.account_id !== 0 && player.name !== "" && player.pattern !== ""
    );
  };
</script>

<Modal isOpen={open} {toggle} on:open={onOpen}>
  <ModalBody class="modal-color">
    <FormGroup>
      <Label>プレイヤー名</Label>
      <Svelecte
        style="color: #2d2c2c;"
        valueAsObject
        id="player"
        fetch={SearchPlayer}
        placeholder=""
        minQuery="3"
        labelField="nickname"
        bind:value={searchResult}
      />
    </FormGroup>

    <FormGroup>
      <Label>アイコン</Label>
      {#await AlertPatterns() then alertPatterns}
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
      {/await}
    </FormGroup>

    <FormGroup>
      <Label>メモ(任意)</Label>
      <Input
        type="textarea"
        maxlength={MAX_MEMO_LENGTH}
        bind:value={target.message}
      />
      <span>{target.message.length}/{MAX_MEMO_LENGTH}</span>
    </FormGroup>
  </ModalBody>

  <ModalFooter class="modal-color">
    <Button size="sm" color="secondary" on:click={toggle}>キャンセル</Button>
    <Button size="sm" color="primary" on:click={() => add(target)}>追加</Button>
  </ModalFooter>
</Modal>
