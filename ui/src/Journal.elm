module Journal exposing (Entry, entriesDecoder, entryDecoder, entryEncoder)

import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Extra exposing (datetime)
import Json.Decode.Pipeline exposing (optional, required)
import Json.Encode as Encode
import Time exposing (Posix)


type alias Entry =
    { id : String
    , content : String
    , created : Posix
    }


entryDecoder : Decoder Entry
entryDecoder =
    Decode.succeed Entry
        |> required "id" string
        |> required "content" string
        |> required "created" datetime


entriesDecoder : Decoder (List Entry)
entriesDecoder =
    Decode.list entryDecoder


entryEncoder : Entry -> Encode.Value
entryEncoder entry =
    Encode.object
        [ ( "content", Encode.string entry.content )
        ]
