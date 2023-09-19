<script lang="ts">
  // @ts-ignore
  import Svelecte from "svelecte";
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
  import { MAX_MEMO_LENGTH } from "src/lib/types";
  import clone from "clone";

  const EMPTY_ALERT_PLAYER: domain.AlertPlayer = {
    account_id: 0,
    name: "",
    pattern: "bi-check-circle-fill",
    message: "",
  };

  export const openAsAdd = () => {
    target = clone(EMPTY_ALERT_PLAYER);
    searchResult = undefined;
    disablePlayerSearch = false;
    isOpen = true;
  };

  export const openAsEdit = (p: domain.AlertPlayer) => {
    target = p;
    searchResult = {
      account_id: p.account_id,
      nickname: p.name,
    };
    disablePlayerSearch = true;
    isOpen = true;
  };

  let isOpen = false;
  let disablePlayerSearch = false;
  let target: domain.AlertPlayer;
  let searchResult: domain.WGAccountListData | undefined;

  const dispatch = createEventDispatcher();

  const update = async (player: domain.AlertPlayer) => {
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
      isOpen = false;
    }
  };

  const validate = (player: domain.AlertPlayer): boolean => {
    return (
      player.account_id !== 0 && player.name !== "" && player.pattern !== ""
    );
  };
</script>

<Modal {isOpen}>
  <ModalBody class="modal-color">
    <FormGroup>
      <Label>プレイヤー名</Label>
      <Svelecte
        valueAsObject
        fetch={SearchPlayer}
        placeholder=""
        minQuery="3"
        labelField="nickname"
        bind:value={searchResult}
        disabled={disablePlayerSearch}
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
    <Button size="sm" color="secondary" on:click={() => (isOpen = false)}
      >キャンセル</Button
    >
    <Button size="sm" color="primary" on:click={() => update(target)}
      >更新</Button
    >
  </ModalFooter>
</Modal>
