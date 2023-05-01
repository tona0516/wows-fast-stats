<script>
import Accordion from "./Accordion.svelte";
import HelpExcludeCheckbox from "./HelpExcludeCheckbox.svelte";
import HelpPr from "./HelpPr.svelte";
</script>

<div class="mt-3 help">
  <Accordion
    id="PrColor"
    title="プレイヤーの配色について"
    component="{HelpPr}"
  />
  <Accordion
    id="ExcludeCheckbox"
    title="プレイヤーのチェックボックスについて"
    component="{HelpExcludeCheckbox}"
  />
</div>
