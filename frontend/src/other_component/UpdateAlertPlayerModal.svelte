<script lang="ts">
  import { MAX_MEMO_LENGTH } from "src/const";

  // @ts-ignore
  import Svelecte from "svelecte";
  import { createEventDispatcher } from "svelte";
  import {
    Button,
    FormGroup,
    Input,
    Label,
    Modal,
    ModalBody,
    ModalFooter,
  } from "sveltestrap";
  import { AlertPatterns, UpdateAlertPlayer } from "wailsjs/go/main/App";
  import type { domain } from "wailsjs/go/models";

  export const toggle = () => (open = !open);
  export const setTarget = (p: domain.AlertPlayer) => (target = p);

  let open = false;
  let target: domain.AlertPlayer;

  const dispatch = createEventDispatcher();

  const update = async (player: domain.AlertPlayer) => {
    try {
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

<Modal isOpen={open} {toggle}>
  <ModalBody class="modal-color">
    <FormGroup>
      <Label>プレイヤー名</Label>
      <Svelecte
        style="color: #2d2c2c;"
        id="player"
        placeholder={target.name}
        disabled={true}
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
    <Button size="sm" color="primary" on:click={() => update(target)}
      >更新</Button
    >
  </ModalFooter>
</Modal>
