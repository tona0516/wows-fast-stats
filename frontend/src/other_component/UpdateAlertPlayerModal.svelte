<script lang="ts">
import Svelecte from "svelecte";
import { createEventDispatcher } from "svelte";
import { Button, Modal, ModalBody, ModalFooter } from "sveltestrap";
import { AlertPatterns, UpdateAlertPlayer } from "../../wailsjs/go/main/App";
import type { vo } from "../../wailsjs/go/models";
import { Const } from "../Const";

export const toggle = () => (open = !open);
export const setTarget = (p: vo.AlertPlayer) => (target = p);

const dispatch = createEventDispatcher();

let open = false;

let target: vo.AlertPlayer;
let alertPatterns: string[] = [];

async function onOpen() {
  alertPatterns = await AlertPatterns();
}

async function update(player: vo.AlertPlayer) {
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
}

function validate(player: vo.AlertPlayer): boolean {
  return player.account_id !== 0 && player.name !== "" && player.pattern !== "";
}
</script>

<Modal isOpen="{open}" toggle="{toggle}" on:open="{onOpen}">
  <ModalBody style="background-color: #2d2c2c;">
    <div class="m-1">
      <label for="player" class="col-form-label">プレイヤー名:</label>

      <Svelecte
        style="color: #2d2c2c;"
        id="player"
        placeholder="{target.name}"
        disabled="{true}"
      />
    </div>

    <div class="m-1">
      <div class="form-group">
        <label for="pattern" class="col-form-label">アイコン:</label>

        <div>
          {#each alertPatterns as pattern}
            <div class="form-check form-check-inline">
              <input
                class="form-check-input"
                type="radio"
                bind:group="{target.pattern}"
                value="{pattern}"
              />
              <label class="form-check-label" for="icon">
                <i class="bi {pattern}"></i>
              </label>
            </div>
          {/each}
        </div>
      </div>
    </div>

    <div class="m-1">
      <label for="message" class="col-form-label">メモ(任意):</label>
      <textarea
        class="form-control"
        id="message"
        maxlength="{Const.MAX_MEMO_LENGTH}"
        bind:value="{target.message}"></textarea>
      <span>{target.message.length}/{Const.MAX_MEMO_LENGTH}</span>
    </div>
  </ModalBody>
  <ModalFooter style="background-color: #2d2c2c;">
    <Button color="secondary" on:click="{toggle}">キャンセル</Button>
    <Button color="primary" on:click="{() => update(target)}">追加</Button>
  </ModalFooter>
</Modal>
