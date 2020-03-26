module Book exposing (Book, bookDecoder, booksDecoder)

import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline exposing (required, optional)

type alias Book =
  { id : String
  , displayName : String
  }

bookDecoder : Decoder Book
bookDecoder =
  Decode.succeed Book
    |> required "id" string
    |> required "display_name" string

booksDecoder : Decoder (List Book)
booksDecoder =
  Decode.list bookDecoder
