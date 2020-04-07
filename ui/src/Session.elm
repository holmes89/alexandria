port module Session exposing (Session(..), Token, storageDecoder, storeCredWith)

import Json.Decode as Decode exposing (Decoder, Value, field, string)
import Json.Encode as Encode


type Session
    = Authenticated Token
    | Unauthenticated


type alias Token =
    String


storeCredWith : Token -> Cmd msg
storeCredWith token =
    let
        json =
            Encode.object
                [ ( "token", Encode.string token )
                ]
    in
    storeCache (Just json)


logout : Cmd msg
logout =
    storeCache Nothing


port storeCache : Maybe Value -> Cmd msg


storageDecoder : Decoder Token
storageDecoder =
    Decode.field "token" string
