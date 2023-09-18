<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import { sampleTeam } from "src/lib/rating/RatingConst";
  import { displayItems } from "src/lib/util";
  import ConfirmModal from "src/other_component/ConfirmModal.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
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
    Table,
  } from "sveltestrap";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;
  export let defaultUserConfig: domain.UserConfig;

  let resetModal: ConfirmModal;

  const dispatch = createEventDispatcher();
  const displays = displayItems();

  const onResetConfirmed = () => {
    inputUserConfig.font_size = defaultUserConfig.font_size;
    inputUserConfig.displays = defaultUserConfig.displays;
    inputUserConfig.custom_color = defaultUserConfig.custom_color;
    inputUserConfig.custom_digit = defaultUserConfig.custom_digit;

    dispatch("Change");
  };
</script>

<ConfirmModal
  bind:this={resetModal}
  message="表示設定をリセットしますか？"
  on:onConfirmed={onResetConfirmed}
/>

<Container fluid class="mt-2">
  <FormGroup>
    <Row>
      <Col>
        <Label>文字サイズ</Label>
      </Col>
    </Row>
    <Row>
      <Col sm="auto">
        <Input
          type="select"
          bind:value={inputUserConfig.font_size}
          on:change={() => dispatch("Change")}
        >
          {#each DispName.FONT_SIZES.toArray() as fs}
            <option
              selected={fs.key === $storedUserConfig.font_size}
              value={fs.key}>{fs.value}</option
            >
          {/each}
        </Input>
      </Col>
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col>
        <Label>表示項目</Label>
      </Col>
    </Row>
    <Row>
      <Col sm="auto">
        <Table class="display-config-table">
          <thead>
            <tr>
              {#each ["項目", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
                <th>{columns}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each displays as display}
              <tr>
                <td>
                  {display.name}
                </td>

                {#if display.shipKey}
                  <td>
                    <Input
                      type="switch"
                      class="center"
                      bind:checked={inputUserConfig.displays.ship[
                        display.shipKey
                      ]}
                      on:change={() => dispatch("Change")}
                    />
                  </td>
                {:else}
                  <td></td>
                {/if}

                {#if display.overallKey}
                  <td>
                    <Input
                      type="switch"
                      class="center"
                      bind:checked={inputUserConfig.displays.overall[
                        display.overallKey
                      ]}
                      on:change={() => dispatch("Change")}
                    />
                  </td>
                {:else}
                  <td></td>
                {/if}

                <td>
                  <Input
                    type="select"
                    bind:value={inputUserConfig.custom_digit[display.digitKey]}
                    on:change={() => dispatch("Change")}
                  >
                    {#each [0, 1, 2] as digit}
                      <option
                        selected={digit ===
                          inputUserConfig.custom_digit[display.digitKey]}
                        value={digit}>{digit}</option
                      >
                    {/each}
                  </Input>
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </Col>
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col>
        <Label>各種カラー</Label>
      </Col>
    </Row>
    <Row>
      <Col sm="auto">
        <Table class="display-config-table">
          <thead>
            <tr>
              {#each ["スキル", "文字色", "背景色"] as column}
                <th>{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.SKILL_LEVELS.toArray() as sl}
              <tr>
                <td>{sl.value}</td>
                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.skill.text[
                        sl.key
                      ]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>

                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.skill.background[
                        sl.key
                      ]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </Col>
      <Col sm="auto">
        <Table class="display-config-table">
          <thead>
            <tr>
              {#each ["Tier", "使用艦", "非使用艦"] as column}
                <th>{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.TIER_GROUPS.toArray() as tg}
              <tr>
                <td>{tg.value}</td>
                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.tier.own[tg.key]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>

                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.tier.other[
                        tg.key
                      ]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </Col>
      <Col sm="auto">
        <Table class="display-config-table">
          <thead>
            <tr>
              {#each ["艦種", "使用艦", "非使用艦"] as column}
                <th>{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.SHIP_TYPES.toArray() as st}
              <tr>
                <td>{st.value}</td>
                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.ship_type.own[
                        st.key
                      ]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>

                <td>
                  <div class="center">
                    <Input
                      type="color"
                      bind:value={inputUserConfig.custom_color.ship_type.other[
                        st.key
                      ]}
                      on:input={() => dispatch("Change")}
                    />
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </Col>
    </Row>
  </FormGroup>

  <FormGroup>
    <Row>
      <Col>
        <Label>プレイヤー名の背景色</Label>
      </Col>
    </Row>
    <Row>
      <Col sm="auto">
        <Input
          type="select"
          bind:value={inputUserConfig.custom_color.player_name}
          on:change={() => dispatch("Change")}
        >
          {#each DispName.PLAYER_NAME_COLORS.toArray() as pnc}
            <option
              selected={pnc.key === $storedUserConfig.custom_color.player_name}
              value={pnc.key}>{pnc.value}</option
            >
          {/each}
        </Input>
      </Col>
    </Row>
  </FormGroup>

  <Row>
    <Col><Label>プレビュー</Label></Col>
  </Row>

  <Row>
    <Col style="font-size: {inputUserConfig.font_size};">
      <StatisticsTable teams={[sampleTeam()]} userConfig={inputUserConfig} />
    </Col>
  </Row>

  <Row>
    <Col sm={{ size: 2, offset: 5 }}>
      <Button color="warning" on:click={resetModal.toggle}>リセット</Button>
    </Col>
  </Row>
</Container>

<style>
  :global(.display-config-table) {
    color: var(--app-text-color);
    text-align: center;
  }
  :global(.display-config-table) th {
    padding-top: 4px;
    padding-bottom: 4px;
  }
  :global(.display-config-table) td {
    padding-top: 4px;
    padding-bottom: 4px;
  }
</style>
