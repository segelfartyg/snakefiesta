import type { Direction } from "../enums/Direction";

export interface PlayerMovementIntent{
    playerId: String,
    timestamp: String,
    direction: Direction

}