// lib/sound.ts

type SoundName =
  | "accepted"
  | "rejected"
  | "beep"
  | "winner"
  | "rules"
  | "roundScores";

class SoundManager {
  private enabled = false;
  private sounds: Record<SoundName, HTMLAudioElement>;

  constructor() {
    this.sounds = {
      accepted: new Audio("/sounds/accepted.wav"),
      rejected: new Audio("/sounds/rejected.wav"),
      beep: new Audio("/sounds/beep.wav"),
      winner: new Audio("/sounds/winner.wav"),
      rules: new Audio("/sounds/rules.wav"),
      roundScores: new Audio("/sounds/score.wav"),
    };

    Object.values(this.sounds).forEach((audio) => {
      audio.volume = 0.5;
    });
  }

  enable() {
    this.enabled = true;
  }

  disable() {
    this.enabled = false;
  }

  toggle() {
    this.enabled = !this.enabled;
    return this.enabled;
  }

  play(name: SoundName) {
    if (!this.enabled) {
      return;
    }
    const sound = this.sounds[name];
    sound.currentTime = 0;
    sound.play().catch(() => {});
  }

  stop(name: SoundName) {
    if (!this.enabled) {
      return;
    }
    const sound = this.sounds[name];
    sound.pause();
    sound.currentTime = 0;
  }
}

export const soundManager = new SoundManager();
