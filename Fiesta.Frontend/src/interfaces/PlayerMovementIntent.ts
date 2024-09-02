import type { Direction } from "../enums/Direction";

export interface PlayerMovementIntent{
    x: Number,
    y: Number,
    playerId: String,
    timestamp: String,
    direction: Direction

}