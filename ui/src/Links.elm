module Links exposing (Link, linkDecoder, linkEncoder, linksDecoder)

import Json.Decode as Decode exposing (Decoder, field, list, string)
import Json.Decode.Extra exposing (datetime)
import Json.Decode.Pipeline exposing (optional, required)
import Json.Encode as Encode
import Time exposing (Posix)


type alias Link =
    { id : String
    , link : String
    , displayName : String
    , iconPath : String
    , created : Posix
    , tags : List String
    }


linkDecoder : Decoder Link
linkDecoder =
    Decode.succeed Link
        |> required "id" string
        |> required "link" string
        |> required "display_name" string
        |> required "icon_path" string
        |> required "created" datetime
        |> required "tag_ids" (list string)


linksDecoder : Decoder (List Link)
linksDecoder =
    Decode.list linkDecoder


linkEncoder : Link -> Encode.Value
linkEncoder link =
    Encode.object
        [ ( "link", Encode.string link.link )
        ]
