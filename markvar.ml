open Yojson.Basic.Util
open Str

(** [read_file filename] reads the contents of a file specified by [filename]
    and returns the content as a list of strings, each string representing a line. *)
let read_file filename : string list =
  let ic = open_in filename in
  let try_read () = try Some (input_line ic) with End_of_file -> None in
  let rec loop acc =
    match try_read () with
    | Some s -> loop (s :: acc)
    | None ->
        close_in ic;
        List.rev acc
  in
  loop []

(** [write_file filename lines] writes a list of strings [lines] to a file specified
    by [filename]. Each string is written as a separate line. *)
let write_file filename lines =
  let oc = open_out filename in
  List.iter (fun line -> Printf.fprintf oc "%s\n" line) lines;
  close_out oc

(** [replace_ids markdown_file json_file] processes a Markdown file specified by
    [markdown_file] and a JSON file specified by [json_file]. It replaces content
    in the Markdown file based on ID-to-content mappings defined in the JSON file. *)
let replace_ids markdown_file json_file =
  let markdown_lines = read_file markdown_file in
  let json = Yojson.Basic.from_file json_file in
  let mappings =
    json |> to_assoc |> List.map (fun (k, v) -> (k, to_string v))
  in
  let replace_tag line =
    let regex = regexp "<!--id:\\(.*?\\)-->.+?<!---->" in
    global_substitute regex
      (fun s ->
        let id = matched_group 1 s in
        try
          let replacement = List.assoc id mappings in
          "<!--id:" ^ id ^ "-->" ^ replacement ^ "<!---->"
        with Not_found -> s)
      line
  in
  List.map replace_tag markdown_lines

(** The main entry point of the program. It expects two command-line arguments:
    the path to a Markdown file and the path to a JSON file. *)
let () =
  let markdown_file = Sys.argv.(1) in
  let json_file = Sys.argv.(2) in
  let updated_lines = replace_ids markdown_file json_file in
  write_file markdown_file updated_lines
