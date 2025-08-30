
<script lang="ts">
  import { ModalManager } from "@libs/ModalManager";
  import { storedAlertPlayers } from "@libs/stores";
</script>


<div class="container mx-auto max-w-3xl py-3 flex flex-col gap-4">
  <div class="card bg-base-100 shadow-xl rounded-xl p-6">
    <div class="flex items-center justify-between mb-4">
      <span class="text-2xl font-bold">アラートプレイヤー管理</span>
      <button
        class="btn btn-primary"
        on:click={() => ModalManager.instance.openForCreate()}
      >追加</button>
    </div>
    <div class="overflow-x-auto rounded-xl shadow border border-base-300 bg-base-200">
      <table class="table w-full text-nowrap">
        <thead>
          <tr class="bg-base-300 text-base-content font-semibold">
            {#each ["ID", "プレイヤー名", "アイコン", "メッセージ", "操作"] as header}
              <th class="text-center">{header}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each $storedAlertPlayers as player, i}
            <tr class="hover:bg-base-100">
              <td class="text-center">{player.account_id}</td>
              <td class="text-center">{player.name}</td>
              <td class="text-center"><i class={`bi ${player.pattern}`}></i></td>
              <td class="text-center">{player.message}</td>
              <td class="text-center">
                <div class="flex justify-center gap-2">
                  <button
                    class="btn btn-sm btn-info"
                    on:click={() => ModalManager.instance.openForEdit(player)}
                  >編集</button>
                  <button
                    class="btn btn-sm btn-error"
                    on:click={() => ModalManager.instance.openForRemove(player)}
                  >削除</button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>
