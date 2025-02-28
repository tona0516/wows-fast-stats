// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {model} from '../models';

export function AlertPatterns():Promise<Array<string>>;

export function AlertPlayers():Promise<Array<model.AlertPlayer>>;

export function AutoScreenshot(arg1:string,arg2:string):Promise<void>;

export function Battle():Promise<model.Battle>;

export function DefaultUserConfig():Promise<model.UserConfigV2>;

export function LatestRelease():Promise<model.LatestRelease>;

export function LogError(arg1:string,arg2:{[key: string]: string}):Promise<void>;

export function LogInfo(arg1:string,arg2:{[key: string]: string}):Promise<void>;

export function ManualScreenshot(arg1:string,arg2:string):Promise<boolean>;

export function MigrateIfNeeded():Promise<void>;

export function OpenDirectory(arg1:string):Promise<void>;

export function RemoveAlertPlayer(arg1:number):Promise<void>;

export function SearchPlayer(arg1:string):Promise<{[key: string]: number}>;

export function SelectDirectory():Promise<string>;

export function Semver():Promise<string>;

export function StartWatching():Promise<void>;

export function UpdateAlertPlayer(arg1:model.AlertPlayer):Promise<void>;

export function UpdateInstallPath(arg1:string):Promise<void>;

export function UpdateUserConfig(arg1:model.UserConfigV2):Promise<void>;

export function UserConfig():Promise<model.UserConfigV2>;

export function ValidateInstallPath(arg1:string):Promise<string>;
