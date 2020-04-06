module Session exposing (Session(..), Token)


type Session
    = Authenticated Token
    | Unauthenticated


type alias Token =
    String
