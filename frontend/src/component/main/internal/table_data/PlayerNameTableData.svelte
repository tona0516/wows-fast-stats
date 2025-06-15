<script lang="ts">
import ExternalLink from "src/component/common/ExternalLink.svelte";
import { ExcludedPlayers } from "src/lib/ExcludedPlayers";
import { Notifier } from "src/lib/Notifier";
import { NumbersURL } from "src/lib/NumbersURL";
import type { PlayerName } from "src/lib/column/model/PlayerName";
import { storedAlertPlayers, storedExcludedPlayers } from "src/stores";
import type { data } from "wailsjs/go/models";
import { ClipboardSetText } from "wailsjs/runtime/runtime";

export let column: PlayerName;
export let player: data.Player;

$: accountID = player.player_info.id;
$: isChecked = !$storedExcludedPlayers.has(accountID);
$: alertPlayer = $storedAlertPlayers.find((it) => it.account_id === accountID);
$: clanTag = column.clanTag(player);
$: isNPC = column.isNPC(player);

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const onCheck = async (e: any) => {
  if (e.target.checked) {
    ExcludedPlayers.remove(accountID);
  } else {
    ExcludedPlayers.add(accountID);
  }
};

const setPlayerNameToClipboard = async () => {
  const isSuccess = await ClipboardSetText(player.player_info.name);

  isSuccess
    ? Notifier.success("コピーしました！")
    : Notifier.failure("コピーに失敗しました");
};
</script>

<td>
  {#if isNPC}
    {column.playerName(player)}
  {:else}
    <div class="flex items-center">
      <input
        class="checkbox"
        type="checkbox"
        on:click={onCheck}
        checked={isChecked}
      />

      <div class="dropdown">
        <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
        <div tabindex="0" class="tooltip" data-tip={alertPlayer?.message}>
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a href="#">
            <div class="flex">
              {#if alertPlayer}
                <i class="bi {alertPlayer.pattern}" />
              {/if}
              {#if clanTag}
                {#if column.clanFlagIconClass(player)}
                  <span class={column.clanFlagIconClass(player)}></span>
                {/if}
                <span style="color: {column.clanColorCode(player)}">
                  {clanTag}
                </span>
              {/if}
              <div style="color: {column.textColorCode(player)}">
                {column.playerName(player)}
              </div>
            </div>
          </a>
        </div>

        <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
        <ul
          tabindex="0"
          class="dropdown-content menu rounded-box w-52 shadow-sm"
        >
          {#if clanTag}
            <li>
              <ExternalLink url={NumbersURL.clan(player.player_info.clan.id)}
                >クラン詳細(WoWS Stats & Numbers)</ExternalLink
              >
            </li>
          {/if}

          <li>
            <ExternalLink
              url={NumbersURL.player(
                player.player_info.id,
                player.player_info.name,
              )}>プレイヤー詳細(WoWS Stats & Numbers)</ExternalLink
            >
          </li>

          <li>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a href="#" on:click={setPlayerNameToClipboard}
              >プレイヤー名をクリップボードにコピーする</a
            >
          </li>
        </ul>
      </div>
    </div>
  {/if}
</td>
