type SoundName =
  | "accepted"
  | "rejected"
  | "beep"
  | "winner"
  | "default"
  | "over";

class SoundManager {
  private enabled = false;
  private sounds: Record<SoundName, HTMLAudioElement>;
  private lastPlayedAt: Partial<Record<SoundName, number>> = {};
  private cooldownMs: Partial<Record<SoundName, number>> = {
    default: 200,
    accepted: 75,
    rejected: 75,
    winner: 250,
    over: 250,
    beep: 0,
  };
  private playInFlight: Partial<Record<SoundName, boolean>> = {};

  constructor() {
    this.sounds = {
      accepted: new Audio("/sounds/accepted.wav"),
      rejected: new Audio("/sounds/rejected.wav"),
      beep: new Audio("/sounds/beep.wav"),
      winner: new Audio("/sounds/winner.wav"),
      default: new Audio("/sounds/default.wav"),
      over: new Audio("/sounds/finished.wav"),
    };

    Object.values(this.sounds).forEach((audio) => {
      audio.volume = 0.5;
      audio.preload = "auto";
    });
  }

  enable() {
    this.enabled = true;
  }

  disable() {
    this.enabled = false;

    Object.values(this.sounds).forEach((audio) => {
      audio.pause();
      audio.currentTime = 0;
    });
  }

  toggle() {
    this.enabled = !this.enabled;

    if (!this.enabled) {
      Object.values(this.sounds).forEach((audio) => {
        audio.pause();
        audio.currentTime = 0;
      });
    }

    return this.enabled;
  }

  play(name: SoundName) {
    if (!this.enabled) {
      return;
    }

    const now = Date.now();
    const cooldown = this.cooldownMs[name] ?? 0;
    const last = this.lastPlayedAt[name] ?? 0;

    if (cooldown > 0 && now - last < cooldown) {
      return;
    }

    if (this.playInFlight[name]) {
      return;
    }

    this.lastPlayedAt[name] = now;

    const sound = this.sounds[name];

    try {
      sound.pause();
      sound.currentTime = 0;

      const maybePromise = sound.play();

      if (maybePromise && typeof maybePromise.then === "function") {
        this.playInFlight[name] = true;
        maybePromise
          .catch(() => {})
          .finally(() => {
            this.playInFlight[name] = false;
          });
      }
    } catch {
      this.playInFlight[name] = false;
    }
  }

  stop(name: SoundName) {
    const sound = this.sounds[name];
    sound.pause();
    sound.currentTime = 0;
    this.playInFlight[name] = false;
  }
}

export const soundManager = new SoundManager();