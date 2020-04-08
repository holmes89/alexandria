module Journal exposing (Entry, entriesDecoder)

import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline exposing (optional, required)


type alias Entry =
    { id : String
    , content : String
    , created : String
    }


entryDecoder : Decoder Entry
entryDecoder =
    Decode.succeed Entry
        |> required "id" string
        |> required "content" string
        |> required "created" string


entriesDecoder : Decoder (List Entry)
entriesDecoder =
    Decode.list entryDecoder
