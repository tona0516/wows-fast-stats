// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {domain} from '../models';
import {vo} from '../models';

export function AddExcludePlayerID(arg1:number):Promise<void>;

export function AlertPatterns():Promise<Array<string>>;

export function AlertPlayers():Promise<Array<domain.AlertPlayer>>;

export function AppVersion():Promise<vo.Version>;

export function ApplyRequiredUserConfig(arg1:string,arg2:string):Promise<vo.ValidatedResult>;

export function ApplyUserConfig(arg1:domain.UserConfig):Promise<void>;

export function AutoScreenshot(arg1:string,arg2:string):Promise<void>;

export function Battle():Promise<domain.Battle>;

export function DefaultUserConfig():Promise<domain.UserConfig>;

export function ExcludePlayerIDs():Promise<Array<number>>;

export function FontSizes():Promise<Array<string>>;

export function LatestRelease():Promise<domain.GHLatestRelease>;

export function LogErrorForFrontend(arg1:string):Promise<void>;

export function LogParam():Promise<vo.LogParam>;

export function ManualScreenshot(arg1:string,arg2:string):Promise<void>;

export function OpenDirectory(arg1:string):Promise<void>;

export function Ready():Promise<void>;

export function RemoveAlertPlayer(arg1:number):Promise<void>;

export function RemoveExcludePlayerID(arg1:number):Promise<void>;

export function SampleTeams():Promise<Array<domain.Team>>;

export function SearchPlayer(arg1:string):Promise<domain.WGAccountList>;

export function SelectDirectory():Promise<string>;

export function StatsPatterns():Promise<Array<string>>;

export function UpdateAlertPlayer(arg1:domain.AlertPlayer):Promise<void>;

export function UserConfig():Promise<domain.UserConfig>;
