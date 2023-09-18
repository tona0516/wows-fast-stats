<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import clone from "clone";
  import { storedAlertPlayers, storedUserConfig } from "src/stores";
  import type { domain } from "wailsjs/go/models";
  import { Alert, Button, Col, Container, Icon, Row, Table } from "sveltestrap";

  const dispatch = createEventDispatcher();

  const alertPlayerColumns = ["プレイヤー名", "アイコン", "メモ", "操作"];

  function onClickAdd() {
    dispatch("AddAlertPlayer");
  }

  function onClickEdit(player: domain.AlertPlayer) {
    const target = clone(player);
    dispatch("UpdateAlertPlayer", { target: target });
  }

  function onClickRemove(player: domain.AlertPlayer) {
    const target = clone(player);
    dispatch("RemoveAlertPlayer", { target: target });
  }
</script>

<Container fluid class="mt-2">
  <Row>
    <Col>
      <Alert color="info" fade={false}>
        <h5 class="alert-heading text-capitalize">プレイヤー検出機能</h5>
        <ul class="m-0">
          <li>戦闘情報のテーブル内のプレイヤーにアイコンを表示</li>
          <li>マウスオーバーでメモを表示</li>
          <li>プレイヤー名をクリックで追加・削除が可能</li>
        </ul>
      </Alert>
    </Col>
  </Row>
  <Row>
    <Col sm={{ size: 8, offset: 2 }}>
      {#if $storedAlertPlayers.length === 0}
        <p>プレイヤーリストがありません</p>
      {:else}
        <div class="table-responsive">
          <Table size="sm" class="alert-player-table">
            <thead>
              <tr>
                {#each alertPlayerColumns as column}
                  <th>{column}</th>
                {/each}
              </tr>
            </thead>
            <tbody>
              {#each $storedAlertPlayers as player}
                <tr>
                  <td>{player.name}</td>
                  <td><i class="bi {player.pattern}" /></td>
                  <td>{player.message}</td>
                  <td>
                    <div class="d-flex justify-content-center">
                      <Button
                        class="alert-player-button"
                        color="success"
                        on:click={() => onClickEdit(player)}
                        >編集
                      </Button>
                      <Button
                        class="alert-player-button"
                        color="danger"
                        on:click={() => onClickRemove(player)}>削除</Button
                      >
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </Table>
        </div>
      {/if}
    </Col>
  </Row>
  <Row>
    <Col sm={{ size: 2, offset: 5 }}>
      <Button color="primary" on:click={onClickAdd}>追加</Button>
    </Col>
  </Row>
</Container>

<style>
  :global(.alert-player-table) {
    color: var(--app-text-color);
    text-align: center;
  }
  :global(.alert-player-table) th {
    padding-top: 4px;
    padding-bottom: 4px;
  }
  :global(.alert-player-table) td {
    padding-top: 4px;
    padding-bottom: 4px;
  }

  :global(.alert-player-button) {
    margin-left: 4px;
    margin-right: 4px;
  }
</style>
