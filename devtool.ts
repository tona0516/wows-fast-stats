import { Command } from "commander";
import fs from "node:fs";
import path from "node:path";
import archiver from "archiver";
import shelljs from "shelljs";
import readlineSync from "readline-sync";

const APP_NAME = "wows-fast-stats";
const SEMVER = "0.15.0";
const BINARY_NAME = `${APP_NAME}-${SEMVER}.exe`;
const DISCORD_WEBHOOK_JSON = "discord_webhook.json";
const FRONTEND_NPM_COMMAND = "npm --prefix ./frontend";
const COMMON_LDFLAGS = {
  "main.AppName": APP_NAME,
  "main.Semver": SEMVER,
} as const;

function getFormattedLDFlags(isDev: boolean): string {
  let discordWebhook: {
    dev: { alert: string; info: string };
    prod: { alert: string; info: string };
  } = { dev: { alert: "", info: "" }, prod: { alert: "", info: "" } };

  try {
    discordWebhook = JSON.parse(fs.readFileSync(DISCORD_WEBHOOK_JSON, "utf8"));
  } catch (error) {
    console.log(`[WARN] Failed to parse ${DISCORD_WEBHOOK_JSON}: ${error}`);
  }

  let flags: { [key: string]: string | boolean };
  if (isDev) {
    flags = {
      ...COMMON_LDFLAGS,
      "main.IsDev": true,
      "main.AlertDiscordWebhookURL": discordWebhook.dev.alert,
      "main.InfoDiscordWebhookURL": discordWebhook.dev.info,
    };
  } else {
    const alertURL = discordWebhook.prod.alert;
    const infoURL = discordWebhook.prod.info;
    if (!alertURL || !infoURL) {
      throw Error("Discord webhook URL not defined");
    }

    flags = {
      ...COMMON_LDFLAGS,
      "main.AlertDiscordWebhookURL": alertURL,
      "main.InfoDiscordWebhookURL": infoURL,
    };
  }

  return Object.entries(flags)
    .map((it) => `-X ${it[0]}=${it[1]}`)
    .join(" ");
}

function exec(command: string) {
  console.log(`$ ${command}`);
  shelljs.exec(command, { env: { ...process.env, FORCE_COLOR: "true" } });
  console.log("");
}

function setup() {
  const goPkgs = [
    "github.com/wailsapp/wails/v2/cmd/wails@latest",
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
    "go.uber.org/mock/mockgen@latest",
    "github.com/fdaines/arch-go@latest",
  ];
  goPkgs.forEach((pkg) => {
    exec(`go install ${pkg}`);
  });
  exec(`${FRONTEND_NPM_COMMAND} ci`);
}

function genmock() {
  exec("go generate ./...");
}

function lint() {
  exec("golangci-lint run");
  exec(`${FRONTEND_NPM_COMMAND} run check`);
  exec(`${FRONTEND_NPM_COMMAND} run lint`);
}

function fmt() {
  exec("golangci-lint run --fix");
  exec("go fmt");
  exec(`${FRONTEND_NPM_COMMAND} run fmt`);
}

function test() {
  exec("arch-go");
  exec("go test -cover ./... -count=1");
  exec(`${FRONTEND_NPM_COMMAND} run test`);
}

function dev() {
  exec(`wails dev -ldflags "${getFormattedLDFlags(true)}"`);
}

function chbtl() {
  const testReplayPath = "./test_install_dir/replays";
  const tempArenaInfoName = "tempArenaInfo.json";

  const files = fs
    .readdirSync(testReplayPath)
    .filter((file) => file.endsWith(".json") && file !== tempArenaInfoName)
    .map((file) => path.join(file));
  files.forEach((file, i) => {
    console.log(i, file);
  });

  const index = readlineSync.question("index? > ");
  fs.copyFileSync(
    path.join(testReplayPath, files[Number(index)]),
    path.join(testReplayPath, tempArenaInfoName)
  );
}

function build() {
  test();
  exec(
    `wails build -ldflags "${getFormattedLDFlags(
      false
    )}" -platform windows/amd64 -o ${BINARY_NAME} -trimpath`
  );
}

async function pkg() {
  build();

  const tempDir = fs.mkdtempSync("");
  try {
    const dstPath = path.join(tempDir, APP_NAME);

    fs.mkdirSync(dstPath);
    fs.copyFileSync(
      path.join("./build/bin/", BINARY_NAME),
      path.join(dstPath, BINARY_NAME)
    );

    const archive = archiver("zip", {});
    const output = fs.createWriteStream(`${APP_NAME}.zip`);
    archive.pipe(output);

    archive.directory(tempDir, false);

    await archive.finalize();
  } finally {
    fs.rmSync(tempDir, { recursive: true, force: true });
  }
}

function clean() {
  [
    "config/",
    "cache/",
    "temp_arena_info/",
    "screenshot/",
    "persistent_data/",
  ].forEach((target) => {
    fs.rmSync(target, { recursive: true, force: true });
  });
}

function main() {
  const program = new Command();

  program
    .command("setup")
    .description("Install dependencies for development.")
    .action(setup);
  program.command("genmock").description("Generate mocks").action(genmock);
  program
    .command("lint")
    .description("Lint for golang and typescript codes.")
    .action(lint);
  program
    .command("fmt")
    .description("Format golang and typescript codes.")
    .action(fmt);
  program.command("test").description("Execute unit tests.").action(test);
  program
    .command("dev")
    .description("Launch the application for developer mode.")
    .action(dev);
  program
    .command("chbtl")
    .description("Replace 'tempArenaInfo.json' in test install directory.")
    .action(chbtl);
  program
    .command("build")
    .description("Make binary for Windows.")
    .action(build);
  program
    .command("pkg")
    .description("Make zip included the binary for Github Release.")
    .action(pkg);
  program
    .command("clean")
    .description("Remove directories and files created by the application.")
    .action(clean);

  program.parse();
}

main();
