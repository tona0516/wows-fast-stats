<script lang="ts">
  import type { Tonako } from "./Tonako";

  export let tonako: Tonako;
  export let message: string;
</script>

<div class="flex w-full h-screen items-center justify-center">
  <div class="flex flex-col items-center">
    <div class="avatar">
      <div class="w-60 rounded-4xl">
        <img src={tonako.getImgPath()} alt="" />
      </div>
    </div>
    <p class="p-4">{message}</p>
  </div>
</div>
