<div align="center">
  <h2 align="center">GO MokuMovies Discord BOT</h2>

  <p align="center">
    A simple discord bot for getting movies and details built with GO and OMDB API
  </p>
</div>

## Notes

Install Go & [taskfile](https://taskfile.dev), coz taskfile make u happy at building :D

## Install

_Tutorial menggunakan bot ini di local dan menggunakan CLI lokal. Pastikan taskfile udah keinstall ya!!_

1. Buat akun OMDB di https://www.omdbapi.com dan copy apikey-nya.
2. Buat akun discord Developer, dan invite ke discord servermu.. ini lumayan makan waktu buatnya.. cari di google yak :P
3. Buat config.json di dir utama dengan isi
   ```sh
   {
    "token"  : "<Kode token discord bot mu>",
    "api_token"  : "<token api ODMB api mu>"
    }
   ```
3. jalankan ini di powershell / sbg. 
   ```sh
   task buildandrun - untuk build dan run sekaligus
   task build - untuk build file main.go
   task run - untuk run file yang sudah di build
   ```
4. Tada, sekarang kamu bisa gunakan !mokuhelp untuk isi list commandnya

## Yang digunakan disini

Use this space to list resources you find helpful and would like to give credit to. I've included a few of my favorites to kick things off!

* [Discordgo](https://github.com/bwmarrin/discordgo)
* [OMBDBAPI](https://www.omdbapi.com)
