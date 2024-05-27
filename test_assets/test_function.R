#' @title name
#' @description The present function create pseudo-bulk matrix from clustering.output file. The output are three files: _bulklog2, which is not normalized, _bulkColumn, which is z-scoere calculated over each column, _bulkRow, which is z-score calculated over each row
#' @param group, a character string. Two options: sudo or docker, depending to which group the user belongs
#' @param scratch.folderDOCKER, a character string indicating the path of the scratch folder inside the docker
#' @param scratch.folderHOST, a character string indicating the path of the scratch folder inside the host. If not running from docker, this is the character string that indicates the path of the scratch.folder
#' @param file, a character string indicating the path of the file, with file name and extension included
#' @param cl, name and path of the file clustering.output previously generated from clustering algorithm from rcasc
#' @param separator, matrix separator, ',', '\\t'
#' @author Lorem Ipsum, Dolor Sit Amet, Lorem Sentence
#'
#' @examples
#' \dontrun{
#'  bulkClusters(group="docker", scratch.folderDOCKER="/sharedFolder/scratch", scratch.folderHOST="/home/user/scratch"
#'  file="/home/user/temp/setA.csv",separator=",",
#'  cl="/home/user/temp/Results/setA/3/setA_clustering.output.csv")
#'  
#'}
#' @export
name <- function(variables) {

}

