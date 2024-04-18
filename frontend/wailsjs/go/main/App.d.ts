// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {data} from '../models';

export function AddExcludePlayerID(arg1:number):Promise<void>;

export function AlertPatterns():Promise<Array<string>>;

export function AlertPlayers():Promise<Array<data.AlertPlayer>>;

export function ApplyRequiredUserConfig(arg1:string,arg2:string):Promise<data.RequiredConfigError>;

export function ApplyUserConfig(arg1:data.UserConfigV2):Promise<void>;

export function AutoScreenshot(arg1:string,arg2:string):Promise<void>;

export function Battle():Promise<data.Battle>;

export function DefaultUserConfig():Promise<data.UserConfigV2>;

export function ExcludePlayerIDs():Promise<Array<number>>;

export function LatestRelease():Promise<data.GHLatestRelease>;

export function LogError(arg1:string,arg2:{[key: string]: string}):Promise<void>;

export function LogInfo(arg1:string,arg2:{[key: string]: string}):Promise<void>;

export function ManualScreenshot(arg1:string,arg2:string):Promise<boolean>;

export function MigrateIfNeeded():Promise<void>;

export function OpenDirectory(arg1:string):Promise<void>;

export function RemoveAlertPlayer(arg1:number):Promise<void>;

export function RemoveExcludePlayerID(arg1:number):Promise<void>;

export function SearchPlayer(arg1:string):Promise<data.WGAccountList>;

export function SelectDirectory():Promise<string>;

export function Semver():Promise<string>;

export function StartWatching():Promise<void>;

export function UpdateAlertPlayer(arg1:data.AlertPlayer):Promise<void>;

export function UserConfig():Promise<data.UserConfigV2>;

export function ValidateRequiredConfig(arg1:string,arg2:string):Promise<data.RequiredConfigError>;
