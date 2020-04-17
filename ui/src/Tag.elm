module Tag exposing (Tag, fetchTags, tagDecoder, tagsDecoder)

import Http
import Json.Decode as Decode exposing (Decoder, field, list, string)
import Json.Decode.Pipeline exposing (optional, required)
import Session exposing (..)


type alias Tag =
    { id : String
    , displayName : String
    , color : String
    }


tagDecoder : Decoder Tag
tagDecoder =
    Decode.succeed Tag
        |> required "id" string
        |> required "display_name" string
        |> required "color" string


tagsDecoder : Decoder (List Tag)
tagsDecoder =
    Decode.list tagDecoder


fetchTags : Token -> (Result Http.Error (List Tag) -> msg) -> Cmd msg
fetchTags token a =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson a tagsDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/tags/"
        }
